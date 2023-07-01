import Layout from "@/components/Layout"
import { NextPageWithLayout } from "./_app"
import { ReactElement } from "react"
import { useRouter } from "next/navigation"


const Page: NextPageWithLayout = () => {
    const { push } = useRouter();
    push("/projects")
    return (
        <></>
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