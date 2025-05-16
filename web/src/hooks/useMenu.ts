import {Menu} from "@/interfaces/contents/MenuContent";
import {useEffect, useState} from "react";

export const useMenu = () : Menu => {
    const [menu, setMenu] = useState<Menu>({
        menus: []
    })

    useEffect(() => {
        fetch(`contents/json/mainMenu.json`)
            .then((response) => {
                return response.json()
            })
            .then((data) => setMenu(data))
    }, []);

    return menu
}