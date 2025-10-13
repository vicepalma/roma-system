import { HTMLAttributes } from 'react'
import { clsx } from 'clsx'
export function CardHeader({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={clsx('p-6 border-b', className)} {...props} />
}
