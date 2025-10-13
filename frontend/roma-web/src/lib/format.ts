export function formatPercent(value: number, decimals = 1, locale = 'es-CL') {
  // Si viene en [0..1], pásalo a porcentaje
  const v = value <= 1 ? value * 100 : value
  return new Intl.NumberFormat(locale, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(v)
}
