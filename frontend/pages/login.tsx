import Layout from "@/components/Layout"
import { NextPageWithLayout, pb } from "./_app"
import { ReactElement, useContext } from "react"
import { useRouter } from "next/navigation"


const Page: NextPageWithLayout = () => {
    const { push } = useRouter();
    const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const data = {
            email: formData.get("email")?.toString(),
            password: formData.get("password")?.toString()
        }
        if (data.email === undefined || data.password === undefined) {
            alert("missing data");
            return
        }
        try {
            await pb.collection('users').authWithPassword(data.email, data.password);
            push("/projects");   
        }
        catch (err) {
            console.error(err)
            alert("Invalid Username or Password") // TODO: use alert box to show error
        }
    }

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
                <div className="m-3">
                    <button type="submit" className="float-right bg-blue-500 text-white p-2 rounded hover:bg-blue-800 transition w-full">Login</button>
                </div>
            </form>
        </div>
      </main>
    )
}


Page.getLayout = function getLayout(page: ReactElement) {
    return (
        <Layout page="home">
            {page}
        </Layout>
    )
}

export default Page