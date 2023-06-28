import Layout from "@/components/Layout"
import { NextPageWithLayout, pb } from "./_app"
import { ReactElement, useContext, useEffect } from "react"
import { useRouter } from "next/navigation"

const Page: NextPageWithLayout = () => {
    const { push } = useRouter();
    useEffect(() => {
        pb.authStore.clear();

        push("/login");
    }, [pb.authStore.token]);
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