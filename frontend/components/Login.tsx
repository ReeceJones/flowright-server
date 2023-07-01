import { pb } from "@/pages/_app"
import { useState } from "react"
import { BsArrowClockwise } from "react-icons/bs"
import Link from "next/link"

export default function Login() {
    const [isLoading, setIsLoading] = useState(false)

    const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        setIsLoading(true)
        const formData = new FormData(event.currentTarget)
        const data = {
            email: formData.get("email")?.toString(),
            password: formData.get("password")?.toString()
        }
        if (data.email === undefined || data.password === undefined) {
            alert("missing data")
            setIsLoading(false)
            return
        }
        try {
            await pb.collection('users').authWithPassword(data.email, data.password);
        }
        catch (err) {
            console.error(err)
            alert("Invalid Username or Password") // TODO: use alert box to show error
        }
        setIsLoading(false)
    }

    // const spinnerClass = isLoading ? " animate-spin" : ""

    const buttonInnerHTML = (isLoading ? <BsArrowClockwise className="animate-spin text-white mx-auto text-2xl"/> : "Login")

    return (
      <main className="p-12">
        <div className="container mx-auto h-12 max-w-md">
            <form onSubmit={onSubmit}>
                <div className="m-3">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Email</label>
                    <input type="email" id="email" name="email" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="email"/>
                </div>
                <div className="m-3">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Password</label>
                    <input type="password" id="password" name="password" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="password"/>
                </div>
                <div className="m-3 w-full text-center">
                    <Link href="/signup" className="text-blue-500 hover:text-blue-800 transition mx-7">Create Account</Link>
                    <Link href="/resetpassword" className="text-blue-500 hover:text-blue-800 transition mx-7">Forgot Password?</Link>
                </div>
                <div className="m-3">
                    <button type="submit" className="bg-blue-500 text-white p-2 rounded hover:bg-blue-800 transition w-full">{buttonInnerHTML}</button>
                </div>
            </form>
        </div>
      </main>
    )
}