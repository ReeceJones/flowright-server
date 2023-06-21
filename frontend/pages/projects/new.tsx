import Layout from "@/components/Layout"
import { NextPageWithLayout } from "../_app"
import { ReactElement, useEffect } from "react"

const Page: NextPageWithLayout = () => {
    useEffect(() => {
        console.log(window.sessionStorage)
    }, [])
    return (
      <main className="p-12">
        <div className="container mx-auto h-12">

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