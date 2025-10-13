import { HTMLAttributes } from 'react'
import { clsx } from 'clsx'
export function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h3 className={clsx('text-lg font-semibold leading-none tracking-tight', className)} {...props} />
}
