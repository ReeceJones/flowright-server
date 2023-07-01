import NavBar from './NavBar'
import { useState, useEffect } from 'react'
import Login from './Login'
import { pb } from '@/pages/_app'

export default function RootLayout({
  children,
  page
}: {
  children: React.ReactNode
  page?: string
}) {
  const [isAuthed, setIsAuthed] = useState(false)

  useEffect(() => {
    setIsAuthed(pb.authStore.isValid)
    pb.authStore.onChange(() => {
      setIsAuthed(pb.authStore.isValid)
    })
  }, [])

  if (isAuthed) {
    return (
        <>
            <NavBar page={page}/>
            <div className="ml-48">
                {children}
            </div>
        </>
    )
  }
  else {
    return (
        <>
            <NavBar page={page}/>
            <div className="ml-48">
                <Login/>
            </div>
        </>
    )
  }
}
