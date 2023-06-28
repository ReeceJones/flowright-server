import { useState, useEffect } from 'react'
import { pb } from '@/pages/_app';
import {BsGear, BsHouse, BsFolder, BsGraphUp, BsPeople, BsPerson, BsBook, BsQuestionLg, BsMegaphone, BsDoorOpen} from 'react-icons/bs'
import NavBarButton from './NavBarButton'

export default function NavBar({page}: {page?: string}) {
    const [name, setName] = useState(undefined)
    const [logoutDisplay, setLogoutDisplay] = useState('none')

    useEffect(() => {
        if (pb.authStore.model) {
            setName(pb.authStore.model.name);
        }
        else {
            setName(undefined);
        }
    }, [pb, pb.authStore.isValid, pb.authStore.model])

    if (name === undefined) {
        return (
            <nav className="bg-slate-100 border shadow-sm w-48 h-screen float-left fixed">
                <ul className="block p-4 h-full">
                    <div className="bottom-0 fixed w-40">
                        <NavBarButton href="/docs" smallPadding={true}><span className="text-sm"><BsBook className="inline-block mr-1 mb-1"/> Docs</span></NavBarButton>
                        <NavBarButton href="/feedback" smallPadding={true}><span className="text-sm"><BsMegaphone className="inline-block mr-1 mb-1"/> Feedback</span></NavBarButton>
                        <NavBarButton href="/support" smallPadding={true}><span className="text-sm"><BsQuestionLg className="inline-block mr-1 mb-1"/> Support</span></NavBarButton>
                        <hr/>
                        <NavBarButton href="/login" selected={page === "profile"}><span><BsPerson className="inline-block mr-1 mb-1"/> Login</span></NavBarButton>
                    </div>
                </ul>
            </nav>
        )
    }
    else {
        return (
            <nav className="bg-slate-100 border shadow-sm w-48 h-screen float-left fixed">
                <ul className="block p-4 h-full">
                    <NavBarButton href="/home" selected={page === "home"}><span><BsHouse className="inline-block mr-1 mb-1"/> Home</span></NavBarButton>
                    <NavBarButton href="/projects" selected={page === "projects"}><span><BsFolder className="inline-block mr-1 mb-1"/> Projects</span></NavBarButton>
                    <NavBarButton href="/metrics" selected={page === "metrics"}><span><BsGraphUp className="inline-block mr-1 mb-1"/> Metrics</span></NavBarButton>
                    <NavBarButton href="/explore" selected={page === "explore"}><span><BsPeople className="inline-block mr-1 mb-1"/> Explore</span></NavBarButton>
                    <NavBarButton href="/settings" selected={page === "settings"}><span><BsGear className="inline-block mr-1 mb-1"/> Settings</span></NavBarButton>
                    <div className="bottom-0 fixed w-40">
                        <NavBarButton href="/docs" smallPadding={true}><span className="text-sm"><BsBook className="inline-block mr-1 mb-1"/> Docs</span></NavBarButton>
                        <NavBarButton href="/feedback" smallPadding={true}><span className="text-sm"><BsMegaphone className="inline-block mr-1 mb-1"/> Feedback</span></NavBarButton>
                        <NavBarButton href="/support" smallPadding={true}><span className="text-sm"><BsQuestionLg className="inline-block mr-1 mb-1"/> Support</span></NavBarButton>
                        <hr/>
                        <div className="transition" onMouseEnter={() => setLogoutDisplay('block')} onMouseLeave={() => setLogoutDisplay('none')}>
                            <div id="logout-button" style={{ display: logoutDisplay }}>
                                <NavBarButton href="/logout" selected={page === "profile"}><span className="text-sm"><BsDoorOpen className="inline-block mr-1 mb-1"/> Logout</span></NavBarButton>
                            </div>
                            <NavBarButton href="/profile" selected={page === "profile"}><span><BsPerson className="inline-block mr-1 mb-1"/> { name }</span></NavBarButton>
                        </div>
                    </div>
                </ul>
            </nav>
        )
    }
}