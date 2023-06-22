import Layout from "@/components/Layout"
import { NextPageWithLayout } from "../../_app"
import { ReactElement } from "react"
import { useRouter } from "next/router"

const Page: NextPageWithLayout = () => {
    const router = useRouter()
    return (
        <main className="p-12">
            <div className="container mx-auto">
                <div className="m-2 mb-6">
                    <h1 className="text-4xl font-bold">{router.query.slug}</h1>
                </div>
                <div className="bg-gray-100 rounded-md p-4 container">
                    <span className="block mb-1">Link your code to this project:</span>
                    <div className="container w-full rounded bg-gray-800 p-3">
                        <code className="text-sm font-mono text-white">
                            $ flowright link {router.query.slug}
                        </code>
                    </div>
                </div>
            </div>
        </main>
    )
}


Page.getLayout = function getLayout(page: ReactElement) {
    return (
        <Layout page="projects">
            {page}
        </Layout>
    )
}

export default Page