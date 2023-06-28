import Layout from "@/components/Layout"
import { NextPageWithLayout } from "./_app"
import { ReactElement } from "react"

const Page: NextPageWithLayout = () => {
    return (
      <main className="p-12">
        <div className="container mx-auto h-12">
            <div className="columns-1 p-3">
                <div className="container p-4 bg-gray-100 mx-auto rounded">

                </div>
            </div>
            <div className="columns-2 p-3">
                <div className="container p-4 bg-gray-100 mx-auto rounded">

                </div>
                <div className="container p-4 bg-gray-100 mx-auto rounded">

                </div>
            </div>
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