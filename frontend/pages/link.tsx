import Layout from "@/components/Layout"
import { NextPageWithLayout, pb } from "./_app"
import { ReactElement, useState, useEffect } from "react"
import { useSearchParams } from "next/navigation"
import { useRouter } from "next/router"


const Page: NextPageWithLayout = () => {
    const { push } = useRouter();
    const [name, setName] = useState("")
    const [confirmed, setConfirmed] = useState(false)

    useEffect(() => {
        setName(pb.authStore.model?.name || "")
    }, [])

    const errorPage = (
        <div className="container mx-auto max-w-lg p-10 text-center">
            <h1 className="text-3xl font-bold">Invalid Challenge</h1>
        </div>
    )

    const handleCancel = async () => {
        push("/projects")
    }

    const handleConfirm = async () => {
        if (challenge === null) {
            return errorPage
        }
        await pb.collection("auth_link_requests").update(challenge, {
            success: true,
            current_client_jwt: pb.authStore.token
        })
        setConfirmed(true)
    }

    const searchParams = useSearchParams()
    if (searchParams === null) {
        return errorPage
    }

    const challenge = searchParams.get("challenge") as string | null
    if (challenge === null || challenge.length !== 15) {
        return errorPage
    }

    if (!confirmed) {
        return (
            <div className="container mx-auto max-w-lg p-10 text-center h-screen">
                <div className="container bg-gray-100 w-full rounded p-2">
                    <span className="text-lg font-bold">{name}</span>
                    <hr className="my-2"/>
                    <span className="text-lg">Are you sure you want to link this client?</span>
                    <div className="columns-2 flex mt-2">
                        <button className="bg-gray-500 text-white p-2 m-2 rounded hover:bg-gray-800 transition w-full" onClick={handleCancel}>Cancel</button>
                        <button className="bg-blue-500 text-white p-2 m-2 rounded hover:bg-blue-800 transition w-full" onClick={handleConfirm}>Authorize</button>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className="container mx-auto max-w-lg p-10 text-center">
            <h1 className="text-3xl font-bold">Successfully linked!</h1>
        </div>
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