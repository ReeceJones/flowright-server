import { BsBoxArrowUpRight } from 'react-icons/bs'
import Project from '@/types/project'

export default function ProjectRow({project}: {project: Project}) {
    let dotColor = project.status === "Alive" ? "bg-green-400" : (project.status === "Dead" ? "bg-red-500" : "bg-orange-400")
    return (
        <tr className="hover:bg-slate-100 transition border-t hover:cursor-pointer">
            <td className="p-4">
                <div className="container w-20">
                    <div className={"w-2.5 h-2.5 rounded-full inline-block mr-2 " + dotColor}></div><span className="inline-block text-sm">{project.status}</span>
                </div>
            </td>
            <td className="p-4">
                <span className="text-sm">1 day ago</span>
            </td>
            <td className="p-4">
                <a className="hover:underline text-blue-500 hover:text-blue-800 whitespace-nowrap" href={ "/projects/settings/" + project.name }><strong>{ project.long_name }</strong></a>
            </td>
            <td className="p-4 px-10 w-full">
                <div className="container max-w-xl overflow-hidden text-ellipsis whitespace-normal line-clamp-1">
                    <span>{ project.description }</span>
                </div>
            </td>
            <td className="p-4">
                <a href="#" className="text-blue-500 hover:text-blue-800 transition float-right flex flex-nowrap">
                    <span className="font-mono text-sm">{ "localhost:9000/reece/" + project.name }</span>
                    <div className="inline-block align-bottom mx-1.5"><BsBoxArrowUpRight className="align-baseline inline-block"/></div>
                </a>
            </td>
        </tr>
    )
}