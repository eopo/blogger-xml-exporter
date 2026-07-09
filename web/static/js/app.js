// Initializes Tom Select for preset dropdown and populates options as DOM elements.
function updatePresetsDropdown(selectEl, presets) {
	if (!selectEl || !presets || !presets.length) return;

	selectEl.innerHTML = '<option value=""></option>';
	presets.forEach((preset, idx) => {
		const opt = document.createElement('option');
		opt.value = String(idx);
		opt.textContent = preset.label;
		selectEl.appendChild(opt);
	});

	if (selectEl.tomselect) {
		selectEl.tomselect.destroy();
	}
	selectEl.tomselect = new TomSelect(selectEl, {
		placeholder: 'Select template…',
		clearable: true,
	});
	selectEl.tomselect.setValue('');
}

// Initializes Tom Select for preset dropdown.
function initPresetsDropdown(selectEl) {
	if (selectEl.tomselect) return;
	selectEl.tomselect = new TomSelect(selectEl, {
		placeholder: 'Select template…',
		clearable: true,
	});
}

// Formats ISO-8601 timestamp to localized date string. Returns empty string for invalid dates.
function formatDate(iso) {
	if (!iso) {
		return "";
	}
	const date = new Date(iso);
	if (Number.isNaN(date.getTime())) {
		return "";
	}
	return date.toLocaleString("en-US", {
		day: "2-digit",
		month: "2-digit",
		year: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	});
}

// Maps form field width (1-12) to Tailwind grid classes.
// On narrow screens, widths > 6 are capped to 6 (col-span-6).
// Specified as literals so Tailwind's class scanner detects them at build time.
const WIDTH_CLASSES = {
	1: "col-span-1 sm:col-span-1",
	2: "col-span-2 sm:col-span-2",
	3: "col-span-3 sm:col-span-3",
	4: "col-span-4 sm:col-span-4",
	5: "col-span-5 sm:col-span-5",
	6: "col-span-6 sm:col-span-6",
	7: "col-span-6 sm:col-span-7",
	8: "col-span-6 sm:col-span-8",
	9: "col-span-6 sm:col-span-9",
	10: "col-span-6 sm:col-span-10",
	11: "col-span-6 sm:col-span-11",
	12: "col-span-6 sm:col-span-12",
};

// Returns width class for a field. Width 0 distributes evenly across row items.
function widthClassFor(width, rowLength) {
	const w = width > 0 ? width : Math.max(1, Math.floor(12 / Math.max(rowLength, 1)));
	return WIDTH_CLASSES[Math.min(w, 12)] || WIDTH_CLASSES[12];
}

// Extracts filename from Content-Disposition header.
function extractFilename(header) {
	if (!header) {
		return null;
	}
	const match = header.match(/filename="?([^";]+)"?/);
	return match ? match[1] : null;
}

// Flattens deeply nested groups beyond MAX_GROUP_DEPTH so they render as fields at current level.
function expandItems(items, allowedGroupDepth) {
	const out = [];
	for (const item of items || []) {
		if (item.hidden) continue;
		if (item.type === "group" && allowedGroupDepth <= 0) {
			out.push(...expandItems(item.items, allowedGroupDepth));
			continue;
		}
		out.push(item);
	}
	return out;
}

// Alpine.js template recursion depth limit. Nested groups beyond this depth are flattened.
const MAX_GROUP_DEPTH = 1;

// Organizes items into rows based on their "row" property.
// Items with the same positive row number appear on the same line.
// Returns guaranteed single-level structure (no nested arrays).
function buildRows(items, allowedGroupDepth = MAX_GROUP_DEPTH) {
	const expanded = expandItems(items, allowedGroupDepth);
	const rowsByKey = new Map();
	const order = [];
	let ownCounter = 0;
	for (const item of expanded) {
		const key = item.row > 0 ? "r" + item.row : "own" + ownCounter++;
		if (!rowsByKey.has(key)) {
			rowsByKey.set(key, []);
			order.push(key);
		}
		rowsByKey.get(key).push(item);
	}
	return order.map((key) => {
		const rowItems = rowsByKey.get(key);
		return rowItems.map((item) => {
			const widthClass = widthClassFor(item.width || 0, rowItems.length);
			if (item.type === "group") {
				return { item, widthClass, isGroup: true, rows: buildRows(item.items, allowedGroupDepth - 1) };
			}
			return { item, widthClass, isGroup: false };
		});
	});
}

// Main Alpine.js app component. Manages form schema, values, post selection, and XML generation.
function app() {
	return {
		formItems: [],
		site: { title: "", heading: "" },
		assets: {},
		values: {},
		presets: {},
		post: {},
		posts: [],
		postsError: false,
		defaultsLoaded: false,
		status: "",
		statusKind: "info",
		postTomSelect: null,

		async init() {
			try {
				await this.loadFormSchema();
				await this.loadDefaults();
				this.initPostPicker();
			} catch (err) {
				this.setStatus(err.message, "error");
			}
		},

		get hasSelectedPost() {
			return this.post && Object.keys(this.post).length > 0;
		},

		// Form sections with resolved rows and presets.
		get effectiveGroups() {
			const sections = [];
			for (const item of this.formItems || []) {
				if (item.hidden) continue;
				if (item.type === "group") {
					sections.push({
						title: item.title || "",
						collapsible: !!item.collapsible,
						collapsed: !!item.collapsed,
						presets: this.presets[item.title] || [],
						rows: buildRows(item.items),
					});
				} else {
					sections.push({
						title: "",
						collapsible: false,
						collapsed: false,
						presets: [],
						rows: [[{ item, widthClass: widthClassFor(0, 1), isGroup: false }]],
					});
				}
			}
			return sections;
		},

		get statusClass() {
			if (this.statusKind === "error") return "status status--error";
			if (this.statusKind === "success") return "status status--success";
			return "status";
		},

		formatDate,

		// Returns combined static and dynamic options for a select field.
		selectOptionsFor(field) {
			const seen = new Set();
			const opts = [];
			for (const o of field.options || []) {
				if (!seen.has(o.value)) {
					seen.add(o.value);
					opts.push({ value: o.value, text: o.label || o.value });
				}
			}
			for (const v of this.values[field.name + "__options"] || []) {
				if (!seen.has(v)) {
					seen.add(v);
					opts.push({ value: v, text: v });
				}
			}
			return opts;
		},

		setStatus(message, kind = "info") {
			this.status = message;
			this.statusKind = kind;
		},

		// Loads form schema and applies theme colors as CSS variables.
		async loadFormSchema() {
			const response = await fetch("/api/form-schema");
			if (!response.ok) {
				throw new Error("failed to load form schema");
			}
			const schema = await response.json();
			this.formItems = schema.items || [];
			this.site = schema.site || { title: "", heading: "" };
			this.assets = schema.assets || {};
			if (this.site.title) {
				document.title = this.site.title;
			}
			if (this.assets.favicon) {
				this.setFavicon(this.assets.favicon);
			}
			if (schema.theme) {
				const root = document.documentElement;
				if (schema.theme.primaryColor) {
					root.style.setProperty("--color-primary", schema.theme.primaryColor);
				}
				if (schema.theme.darkColor) {
					root.style.setProperty("--color-primary-dark", schema.theme.darkColor);
				}
				if (schema.theme.lightColor) {
					root.style.setProperty("--color-primary-light", schema.theme.lightColor);
				}
			}
		},

		// Loads default form values (without selected post).
		async loadDefaults() {
			const response = await fetch("/api/defaults");
			if (!response.ok) {
				throw new Error("failed to load defaults");
			}
			const data = await response.json();
			this.post = data.post || {};
			this.values = data.values || {};
			this.presets = data.presets || {};
			this.defaultsLoaded = true;
		},

		// Applies a preset: merges preset values into current form values.
		applyPreset(preset) {
			for (const [name, value] of Object.entries(preset.values || {})) {
				this.values[name] = value;
			}
		},

		// Resets fields touched by any preset in the group.
		resetPreset(presets) {
			const fieldsToReset = new Set();
			for (const preset of presets || []) {
				for (const name of Object.keys(preset.values || {})) {
					fieldsToReset.add(name);
				}
			}
			for (const name of fieldsToReset) {
				this.values[name] = '';
			}
		},

		// Initializes Tom Select dropdown with live Blogger API search.
		initPostPicker() {
			this.postTomSelect = new TomSelect(this.$refs.postSelect, {
				valueField: "value",
				labelField: "title",
				searchField: ["title", "dateText"],
				options: [],
				placeholder: "Select post…",
				load: (query, callback) => this.searchPosts(query, callback),
				render: {
					option: (data, escape) =>
						`<div class="post-option"><span class="post-option__title">${escape(data.title)}</span><span class="post-option__date">${escape(data.dateText)}</span></div>`,
					item: (data, escape) => `<div>${escape(data.title)}</div>`,
					no_results: () => `<div class="no-results">No posts found</div>`,
				},
				onChange: (value) => {
					if (value) {
						this.selectPost(value);
					}
				},
			});
			this.postTomSelect.load("");
		},

		// Loads blog posts: recent posts without query, or search results with query.
		async searchPosts(query, callback) {
			this.postsError = false;
			try {
				const url = query ? `/api/posts?q=${encodeURIComponent(query)}` : "/api/posts";
				const response = await fetch(url);
				if (!response.ok) {
					throw new Error("failed to load posts");
				}
				this.posts = await response.json();
				callback(this.posts.map((p) => ({ value: p.id, title: p.title || p.id, dateText: formatDate(p.published) })));
			} catch (err) {
				this.postsError = true;
				callback();
			}
		},

		// Retries post search.
		loadPosts() {
			this.postTomSelect.load(this.postTomSelect.lastValue || "");
		},

		// Clears selected post and reloads default form values.
		async clearPost() {
			this.postTomSelect?.clear(true);
			try {
				await this.loadDefaults();
				this.setStatus("");
			} catch (err) {
				this.setStatus(err.message, "error");
			}
		},

		// Loads a post and its pre-filled values.
		async selectPost(postId) {
			try {
				this.setStatus("Loading post…");
				const response = await fetch(`/api/posts/${encodeURIComponent(postId)}`);
				if (!response.ok) {
					throw new Error("failed to load post");
				}
				const data = await response.json();
				this.post = data.post || {};
				this.values = data.values || {};
				this.presets = data.presets || {};
				this.setStatus("");
			} catch (err) {
				this.setStatus(err.message, "error");
			}
		},

		// Adds a new empty row to an array field.
		addArrayRow(field) {
			if (!Array.isArray(this.values[field.name])) {
				this.values[field.name] = [];
			}
			const row = {};
			for (const sub of field.fields) {
				row[sub.name] = "";
			}
			this.values[field.name].push(row);
		},

		// Removes a row from an array field.
		removeArrayRow(field, idx) {
			this.values[field.name]?.splice(idx, 1);
		},

		// Generates XML and triggers download.
		async submitGenerate() {
			try {
				this.setStatus("Generating XML…");
				const response = await fetch("/api/generate", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ post: this.post, values: this.values }),
				});
				if (!response.ok) {
					const err = await response.json().catch(() => ({}));
					throw new Error(err.error || "failed to generate XML");
				}

				const blob = await response.blob();
				const url = URL.createObjectURL(blob);
				const a = document.createElement("a");
				a.href = url;
				a.download = extractFilename(response.headers.get("Content-Disposition")) || "post.xml";
				document.body.appendChild(a);
				a.click();
				a.remove();
				URL.revokeObjectURL(url);

				this.setStatus("XML downloaded.", "success");
			} catch (err) {
				this.setStatus(err.message, "error");
			}
		},
	};
}
