import NavBar from '../components/NavBar'

export default function RootLayout({
  children,
  page
}: {
  children: React.ReactNode
  page?: string
}) {
  return (
    <>
        <NavBar page={page}/>
        <div className="ml-48">
            {children}
        </div>
    </>
  )
}
