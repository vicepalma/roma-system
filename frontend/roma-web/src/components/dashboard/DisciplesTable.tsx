import { useNavigate } from 'react-router-dom'
import type { CoachDisciple  } from '@/types/coach'

export default function DisciplesTable({ data }: { data: CoachDisciple[] }) {
  const navigate = useNavigate()
  const rows = Array.isArray(data) ? data : []

  return (
    <div className="overflow-x-auto rounded-lg border bg-white">
      <table className="min-w-full text-sm">
        <thead className="bg-gray-50 text-gray-600">
          <tr>
            <th className="text-left px-4 py-2">Nombre</th>
            <th className="text-left px-4 py-2">Email</th>
            <th className="text-left px-4 py-2">Estado</th>
            <th className="text-left px-4 py-2">Última actividad</th>
            <th className="px-4 py-2"></th>
          </tr>
        </thead>
        <tbody>
          {rows.map((d: CoachDisciple) => (
            <tr key={d.id} className="border-t">
              <td className="px-4 py-2">{d.name}</td>
              <td className="px-4 py-2">{d.email || '-'}</td>
              <td className="px-4 py-2 text-right">
                <button
                  onClick={() =>
                    navigate(`/disciples/${d.id}`, {
                      state: { name: d.name, email: d.email },
                    })
                  }
                  className="text-blue-600 hover:underline"
                >
                  Ver detalle
                </button>
              </td>
            </tr>
          ))}
          {rows.length === 0 && (
            <tr><td colSpan={5} className="px-4 py-6 text-center text-gray-500">Sin discípulos aún</td></tr>
          )}
        </tbody>
      </table>
    </div>
  )
}
