import {BsGear, BsHouse, BsFolder, BsGraphUp, BsPeople, BsPerson, BsBook, BsQuestionLg, BsMegaphone} from 'react-icons/bs'
import NavBarButton from './NavBarButton'

export default function NavBar() {
    return (
        <nav className="bg-slate-100 border shadow-sm w-48 h-screen float-left fixed">
            <ul className="block p-4 h-full">
                <NavBarButton href="/home"><span><BsHouse className="inline-block mr-1 mb-1"/> Home</span></NavBarButton>
                <NavBarButton href="/projects" selected={true}><span><BsFolder className="inline-block mr-1 mb-1"/> Projects</span></NavBarButton>
                <NavBarButton href="/metrics"><span><BsGraphUp className="inline-block mr-1 mb-1"/> Metrics</span></NavBarButton>
                <NavBarButton href="/explore"><span><BsPeople className="inline-block mr-1 mb-1"/> Explore</span></NavBarButton>
                <NavBarButton href="/settings"><span><BsGear className="inline-block mr-1 mb-1"/> Settings</span></NavBarButton>
                <div className="bottom-0 fixed w-40">
                    <NavBarButton href="/docs" smallPadding={true}><span className="text-sm"><BsBook className="inline-block mr-1 mb-1"/> Docs</span></NavBarButton>
                    <NavBarButton href="/feedback" smallPadding={true}><span className="text-sm"><BsMegaphone className="inline-block mr-1 mb-1"/> Feedback</span></NavBarButton>
                    <NavBarButton href="/support" smallPadding={true}><span className="text-sm"><BsQuestionLg className="inline-block mr-1 mb-1"/> Support</span></NavBarButton>
                    <hr/>
                    <NavBarButton href="/profile"><span><BsPerson className="inline-block mr-1 mb-1"/> Reece Jones</span></NavBarButton>
                </div>
            </ul>
        </nav>
    )
}