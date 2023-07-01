import { BsBoxArrowUpRight } from 'react-icons/bs'
import Project from '@/types/project'
import Link from 'next/link'

export default function ProjectRow({project}: {project: Project}) {
    const getURL = () => {
        return "localhost:9000/" + project.expand.owner.username + "/" + project.name
    }
    let dotColor = project.status === "Alive" ? "bg-green-400" : (project.status === "Dead" ? "bg-red-500" : "bg-orange-400")
    return (
        <tr className="hover:bg-slate-100 transition border-t">
            <td className="p-4">
                <div className="container w-20">
                    <div className={"w-2.5 h-2.5 rounded-full inline-block mr-2 " + dotColor}></div><span className="inline-block text-sm">{project.status}</span>
                </div>
            </td>
            <td className="p-4">
                <span className="text-sm">{new Date(project.updated).toLocaleDateString('en-us', {year: "numeric", month: "short", day: "numeric"})}</span>
            </td>
            <td className="p-4">
                <Link className="hover:underline text-blue-500 hover:text-blue-800 whitespace-nowrap" href={ "/projects/settings/" + project.name }><strong>{ project.name }</strong></Link>
            </td>
            <td className="p-4 px-10 w-full">
                <div className="container max-w-xl overflow-hidden text-ellipsis whitespace-normal line-clamp-1">
                    <span>{ project.description }</span>
                </div>
            </td>
            <td className="p-4">
                <Link href={"http://" + getURL()} className="text-blue-500 hover:text-blue-800 transition float-right flex flex-nowrap">
                    <span className="font-mono text-sm">{ getURL() }</span>
                    <div className="inline-block align-bottom mx-1.5"><BsBoxArrowUpRight className="align-baseline inline-block"/></div>
                </Link>
            </td>
        </tr>
    )
}