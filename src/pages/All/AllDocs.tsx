import {
    List,
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    Logo,
    useTitle,
} from "@vertex-center/components";
import styles from "./AllDocs.module.sass";
import { useNavigate } from "react-router-dom";

export default function AllDocs() {
    useTitle("All documentation");

    const navigate = useNavigate();

    return (
        <div className={styles.all}>
            <List>
                <ListItem onClick={() => navigate("/vertex")}>
                    <ListIcon>
                        <Logo />
                    </ListIcon>
                    <ListInfo>
                        <ListTitle>Vertex</ListTitle>
                        <ListDescription>
                            The main Vertex documentation.
                        </ListDescription>
                    </ListInfo>
                </ListItem>
            </List>
        </div>
    );
}
