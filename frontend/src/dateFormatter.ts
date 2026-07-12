export function formatDate(val: string | null | undefined, includeTime: boolean | undefined): string {
  if (!val) return ''
  try {
    const d = new Date(val)
    if (isNaN(d.getTime())) return val
    if (includeTime) {
      // browser timezone
      return d.toLocaleDateString(undefined, { 
        year: 'numeric', month: '2-digit', day: '2-digit', 
        hour: '2-digit', minute: '2-digit' 
      })
    }
    return d.toLocaleDateString(undefined, {
      year: 'numeric', month: '2-digit', day: '2-digit'
    })
  } catch {
    return val
  }
}
