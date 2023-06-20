export default function ProjectRowHeader() {
    return (
        <tr className="text-gray-600">
            <th className="p-4 pr-8 float-left"><span className="float-left">Status</span></th>
            <th className="p-4 px-0"><span className="float-left w-32">Last Updated</span></th>
            <th className="p-4"><span className="float-left">Project</span></th>
            <th className="p-4 px-10"><span className="float-left">Description</span></th>
            {/* <th>Created</th> */}
            <th className="p-4"><span className="float-right">Link</span></th>
        </tr>
    )
}