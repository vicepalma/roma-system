export default function ErrorBoundary() {
  return (
    <div className="mx-auto max-w-xl p-6">
      <h2 className="text-lg font-semibold mb-2">Ocurrió un error</h2>
      <p className="text-sm text-gray-600">Intenta recargar la página o volver al dashboard.</p>
      <a href="/dashboard" className="mt-3 inline-block text-blue-600 hover:underline">Ir al dashboard</a>
    </div>
  )
}
