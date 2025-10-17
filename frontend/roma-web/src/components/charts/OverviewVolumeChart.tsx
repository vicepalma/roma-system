import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend } from 'recharts'
import { pivotToChartData, type Pivot as PivotLib } from '@/lib/pivot'
import type { Overview } from '@/types/disciples'

type Props = { overview?: Overview }

export default function OverviewVolumeChart({ overview }: Props) {
  const p = overview?.pivot

  const adapted: PivotLib | null = p
    ? {
        rows: p.data ?? [],
        columns: p.series ?? [],
      }
    : null

  const { data, seriesKeys } = adapted
    ? pivotToChartData(adapted, p?.days ?? 14)
    : { data: [], seriesKeys: [] }

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
