import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts'
import { pivotToChartData } from '@/lib/pivot'
import type { Overview } from '@/types/disciples'

export default function OverviewVolumeChart({ overview }: { overview?: Overview }) {
  const { data, seriesKeys } = pivotToChartData(overview?.pivot, overview?.pivot?.days ?? 14)

  if (!data.length || !seriesKeys.length) {
    return <div className="text-sm text-gray-500">Sin datos suficientes para graficar</div>
  }

  return (
    <div className="h-72 w-full">
      <ResponsiveContainer>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="date" fontSize={12} />
          <YAxis fontSize={12} />
          <Tooltip />
          <Legend />
          {seriesKeys.map((key, idx) => (
            <Bar key={key} dataKey={key} stackId="vol" />
          ))}
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
