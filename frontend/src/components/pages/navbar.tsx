import * as React from "react";
import { ChevronFirst } from "lucide-react";

const SideBar = ({ children }: React.ComponentProps<"aside">) =>
    (
        <aside className="h-screen">
            <div className="h-full flex flex-col bg-white border-r shadow-sm">
                <div className="p-4 pb-2 flex justify-between items-center">
                    <button className="p-1.5 rounded-lg bg-gray-50 hover:bg-gray-100">
                        <ChevronFirst />
                    </button>
                </div>
                <ul className="flex-1 px-3">{ children }</ul>
                <div className="border-t flex p-3">
                    <div className="flex justify-between items-center w-52 ml-3">
                        test
                    </div>
                </div>
            </div>
        </aside>
    )

export {
    SideBar,
}
