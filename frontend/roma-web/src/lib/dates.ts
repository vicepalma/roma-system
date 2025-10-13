export function fmtDate(d?: string | null, locale = 'es-CL') {
  if (!d) return '-'
  const dt = new Date(d + (d.length <= 10 ? 'T00:00:00' : ''))
  if (isNaN(+dt)) return '-'
  return new Intl.DateTimeFormat(locale, { dateStyle: 'medium' }).format(dt)
}

export function fmtDateTime(d?: string | null, locale = 'es-CL') {
  if (!d) return '-'
  const dt = new Date(d)
  if (isNaN(+dt)) return '-'
  return new Intl.DateTimeFormat(locale, { dateStyle: 'medium', timeStyle: 'short' }).format(dt)
}
