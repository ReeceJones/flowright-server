import Layout from "@/components/Layout"
import { NextPageWithLayout } from "../_app"
import { ReactElement, useEffect } from "react"
import PocketBase from "pocketbase"

const pb = new PocketBase("http://localhost:8090")

const Page: NextPageWithLayout = () => {
    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        const formData = new FormData(event.currentTarget)
        const data = {
            long_name: formData.get("name"),
            name: formData.get("name")?.toString().replace(/ /g, "_").replace(/[^a-zA-Z0-9_]/g, "").toLowerCase(),
            description: formData.get("description"),
            // visibility: formData.get("visibility"),
            status: "Alive",
            owner: pb.authStore.model?.id
        }
        console.log(data)
        const res = await pb.collection("projects").create(data)
        console.log(res)
    }
    return (
      <main className="p-12">
        <div className="container mx-auto h-12 ">
            <form onSubmit={handleSubmit}>
                <div className="bg-gray-100 p-6 rounded max-w-3xl container mx-auto">
                    <label className="block text-gray-700 font-bold mb-2" htmlFor="name">Project name*</label>
                    <input className="block w-full shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline mb-4" maxLength={120} id="name" name="name" type="text" placeholder="Project name"/>

                    <label className="block text-gray-700 font-bold mb-2" htmlFor="description">Description*</label>
                    <textarea className="block resize-none w-full h-20 shadow appearance-none border rounded py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline mb-4" maxLength={200} id="description" name="description" placeholder="Description"></textarea>

                    <input className="inline-block" type="radio" id="public" name="visibility" value="public" defaultChecked/>
                    <label className="ml-2 text-lg font-bold" htmlFor="public">
                        <span className="inline-block">Public</span>
                        <span className="block text-sm font-light">Anyone can visit your project site and your project will be visible under the explore tab.</span>
                    </label>

                    <input type="radio" id="private" name="visibility" value="private"/>
                    <label className="ml-2 text-lg font-bold" htmlFor="private">
                        <span className="inline-block">Private</span>
                        <span className="block text-sm font-light">Only you visit your project site.</span>
                    </label>

                    <br/>

                    <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" type="submit">
                        Create Project
                    </button>
                </div>
            </form>
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