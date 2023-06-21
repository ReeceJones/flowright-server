import Layout from "@/components/Layout"
import { NextPageWithLayout } from "./_app"
import { ReactElement, useEffect } from "react"
import ProjectRow from "@/components/ProjectRow"
import ProjectRowHeader from "@/components/ProjectRowHeader"

const Page: NextPageWithLayout = () => {
    useEffect(() => {
        console.log(window.sessionStorage)
    }, [])
    return (
      <main className="p-12">
        <div className="container mx-auto h-12">
          <input type="text" id="fname" name="fname" className="border rounded float-left p-2 hover:bg-slate-100 transition w-96" placeholder="Search..."/>
  
          <a href="/projects/new"><button className="float-right bg-blue-500 text-white p-2 rounded hover:bg-blue-800 transition">New Project</button></a>
        </div>
        <br></br>
        <div className="container mx-auto">
          <div className="border rounded-md">
            <table className="table-auto table">
              <thead className="table-header-group bg-slate-100">
                <ProjectRowHeader/>
              </thead>
              <tbody className="table-row-group">
                <ProjectRow status="Alive" name="My Project" url="localhost:9000/reece/test/"/>
                <ProjectRow status="Pending" name="Resource Calculator" url="localhost:9000/reece/k8s_resource/"/>
                <ProjectRow status="Dead" name="Cool toolz" url="localhost:9000/reece/tools/"/>
                <ProjectRow status="Alive" name="Cloud Resource Monitor" url="localhost:9000/reece/resource_monitor/"/>
                <ProjectRow status="Alive" name="Instance Managaement" url="localhost:9000/reece/instances/"/>
              </tbody>
            </table>
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