import Layout from "@/components/Layout"
import { NextPageWithLayout, pb } from "./_app"
import { ReactElement } from "react"
import { useRouter } from "next/navigation"

const Page: NextPageWithLayout = () => {
    const { push } = useRouter();

    const checkUsername = async (username: string) => {
        const res = await fetch("http://localhost:8090/api/checkusername", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                username: username
            })
        })
        return res
    }

    const checkEmail = async (email: string) => {
        const res = await fetch("http://localhost:8090/api/checkemail", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: email
            })
        })
        return res
    }

    const handleUsernameInput = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const usernameRegex = /^[a-zA-Z0-9_]+$/
        let res = undefined

        if (!(usernameRegex.test(e.target.value))) {
            e.target.setCustomValidity("Invalid username! Only letters, numbers, and underscores are allowed.");
        }
        else if (!(res = await checkUsername(e.target.value)).ok) {
            e.target.setCustomValidity("Invalid username! " + (await res.text()));
        }
        else {
            e.target.setCustomValidity("");
        }

        e.target.reportValidity();
    }

    const handleEmailInput = async (e: React.ChangeEvent<HTMLInputElement>) => {
        let res = undefined;
        if (!(res = await checkEmail(e.target.value)).ok) {
            e.target.setCustomValidity("Invalid email! " + (await res.text()));
        }
        else {
            e.target.setCustomValidity("");
        }
        e.target.reportValidity();
    }

    const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const password = document.getElementById("password") as HTMLInputElement
        const passwordConfirm = document.getElementById("passwordConfirm") as HTMLInputElement
        if (password.value !== passwordConfirm.value) {
            passwordConfirm.setCustomValidity("Passwords do not match!")
        }
        else {
            passwordConfirm.setCustomValidity("")
        }
        passwordConfirm.reportValidity()
    }

    const handleFormSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        if (e.currentTarget.checkValidity() === false) {
            return
        }
        
        const formData = new FormData(e.currentTarget)
        const data = {
            name: formData.get("fname")?.toString(),
            email: formData.get("email")?.toString(),
            username: formData.get("username")?.toString(),
            password: formData.get("password")?.toString(),
            passwordConfirm: formData.get("passwordConfirm")?.toString()
        }
        
        if (data.name === undefined) {
            let formElement = document.getElementById("fname") as HTMLInputElement
            formElement.setCustomValidity("Name is required!")
            formElement.reportValidity()
            return
        }
        if (data.email === undefined) {
            let formElement = document.getElementById("email") as HTMLInputElement
            formElement.setCustomValidity("Email is required!")
            formElement.reportValidity()
            return
        }
        if (data.username === undefined) {
            let formElement = document.getElementById("username") as HTMLInputElement
            formElement.setCustomValidity("Username is required!")
            formElement.reportValidity()
            return
        }
        if (data.password === undefined) {
            let formElement = document.getElementById("password") as HTMLInputElement
            formElement.setCustomValidity("Password is required!")
            formElement.reportValidity()
            return
        }
        if (data.passwordConfirm === undefined) {
            let formElement = document.getElementById("passwordConfirm") as HTMLInputElement
            formElement.setCustomValidity("Password confirmation is required!")
            formElement.reportValidity()
            return
        }

        try {
            await pb.collection("users").create(data)
            await pb.collection("users").authWithPassword(data.email, data.password)
            push("/projects")
        }
        catch (err) {
            console.error(err)
            alert(err)
        }
    }

    return (
      <main className="p-12">
        <div className="container mx-auto h-12 max-w-xl">
            <form onSubmit={handleFormSubmit}>
                <div className="m-3">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Name</label>
                    <input type="text" id="fname" name="fname" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="John Doe"/>
                </div>
                <div className="m-3">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Username</label>
                    <input type="text" id="username" name="username" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="username" onInput={handleUsernameInput}/>
                </div>
                <div className="m-3">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Email</label>
                    <input type="email" id="email" name="email" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="email" onInput={handleEmailInput}/>
                </div>
                <div className="columns-2 m-3">
                    <div className="">
                        <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Password</label>
                        <input type="password" id="password" name="password" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="password" onBlur={handlePasswordChange}/>
                    </div>
                    <div className="">
                        <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Confirm Password</label>
                        <input type="password" id="passwordConfirm" name="passwordConfirm" className="border rounded p-2 hover:bg-slate-100 transition w-full" placeholder="password" onInput={handlePasswordChange}/>
                    </div>
                </div>
                <div className="m-3">
                    <button type="submit" className="float-right bg-blue-500 text-white p-2 rounded hover:bg-blue-800 transition w-full">Create Account</button>
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