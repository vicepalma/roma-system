import { HTMLAttributes } from 'react'
import { clsx } from 'clsx'
export function CardContent({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={clsx('p-6', className)} {...props} />
}
