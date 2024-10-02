import * as React from "react";
import { Link, useLocation } from "react-router-dom"

import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuViewport,
  navigationMenuTriggerStyle
} from "@/components/ui/navigation-menu"

import ThemeToggle from './theme-toggle'


const NavBar: React.FC = () => {
    const location = useLocation()
    if (location.pathname === "/login" || location.pathname.startsWith("/404")) {
        return null
    }

    return (
            <NavigationMenu orientation="vertical">
                <NavigationMenuList className="flex-col items-start">
                    <NavigationMenuItem>
                        <Link to="/">
                            <NavigationMenuLink className={navigationMenuTriggerStyle()}>Home</NavigationMenuLink>
                        </Link>
                    </NavigationMenuItem>

                    <NavigationMenuItem>
                        <Link to="/brothers">
                        <NavigationMenuLink className={navigationMenuTriggerStyle()}>Members</NavigationMenuLink>
                        </Link>
                    </NavigationMenuItem>

                    <NavigationMenuItem>
                        <Link to="/events">
                        <NavigationMenuLink className={navigationMenuTriggerStyle()}>Events</NavigationMenuLink>
                        </Link>
                    </NavigationMenuItem>

                    <NavigationMenuItem>
                        <ThemeToggle />
                    </NavigationMenuItem>

                </NavigationMenuList>
            </NavigationMenu>
    )
}

export default NavBar;
