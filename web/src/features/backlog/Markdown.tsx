import ReactMarkdown from 'react-markdown'

type Props = {
  texto: string
}

export function Markdown({ texto }: Props) {
  return (
    <div className="markdown">
      <ReactMarkdown>{texto}</ReactMarkdown>
    </div>
  )
}
