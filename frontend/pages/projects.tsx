import Layout from "@/components/Layout"
import { NextPageWithLayout } from "./_app"
import { ReactElement, useEffect, useState } from "react"
import Link from "next/link"
import ProjectRow from "@/components/ProjectRow"
import ProjectRowHeader from "@/components/ProjectRowHeader"
import PocketBase, { ListResult, Record } from "pocketbase"
import Project from "@/types/project"

const pb = new PocketBase("http://localhost:8090")

const Page: NextPageWithLayout = () => {
    const projectRows: ReactElement[] = []
    const [projects, setProjects] = useState<Project[]>([])
    useEffect(() => {
        const getProjects = async () => {
            const authData = await pb.collection('users').authWithPassword('reece', 'password');
            const pbProjects = await pb.collection("projects").getList<Project>();
            setProjects(pbProjects.items);
        }
        getProjects();
    }, []);

    for (const project of projects) {
        projectRows.push(<ProjectRow project={project}/>);
    }

    return (
      <main className="p-12">
        <div className="container mx-auto h-12">
          <input type="text" id="fname" name="fname" className="border rounded float-left p-2 hover:bg-slate-100 transition w-96" placeholder="Search..."/>
  
          <Link href="/projects/new"><button className="float-right bg-blue-500 text-white p-2 rounded hover:bg-blue-800 transition">New Project</button></Link>
        </div>
        <br></br>
        <div className="container mx-auto">
          <div className="border rounded-md">
            <table className="table-auto table">
              <thead className="table-header-group bg-slate-100">
                <ProjectRowHeader/>
              </thead>
              <tbody className="table-row-group">
                {/* <ProjectRow status="Alive" name="My Project" url="localhost:9000/reece/test/"/>
                <ProjectRow status="Pending" name="Resource Calculator" url="localhost:9000/reece/k8s_resource/"/>
                <ProjectRow status="Dead" name="Cool toolz" url="localhost:9000/reece/tools/"/>
                <ProjectRow status="Alive" name="Cloud Resource Monitor" url="localhost:9000/reece/resource_monitor/"/>
                <ProjectRow status="Alive" name="Instance Managaement" url="localhost:9000/reece/instances/"/> */}
                {projectRows}
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