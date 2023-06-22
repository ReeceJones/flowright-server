import Link from "next/link"

export default function NavBarButton({children, selected, href, smallPadding}: {children: React.ReactNode, selected?: boolean, href: string, smallPadding?: boolean}) {
    let classes = "my-2 p-2 rounded hover:bg-blue-400 active:bg-blue-500 transition"
    if (selected !== undefined && selected === true) {
        classes += " bg-blue-200"
    }
    if (smallPadding !== undefined && smallPadding === true) {
        classes += " py-1"
    }

    return (
        <li className={classes}>
            <Link href={href}>
                <button className="w-full text-start" >
                    {children}
                </button>
            </Link>
        </li>
    )
}