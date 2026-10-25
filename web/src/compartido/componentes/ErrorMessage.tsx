import type { HTMLAttributes } from 'react'

type ErrorMessageProps = HTMLAttributes<HTMLParagraphElement> & {
  message: string
}

export function ErrorMessage({ message, ...props }: ErrorMessageProps) {
  return <p className="message message--error" role="alert" {...props}>{message}</p>
}
