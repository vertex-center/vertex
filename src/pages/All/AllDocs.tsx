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

const docs = import.meta.glob("/docs/*/doc.json", {
    eager: true,
});

export default function AllDocs() {
    useTitle("All documentation");

    const navigate = useNavigate();

    return (
        <div className={styles.all}>
            <List>
                {
                    Object.entries(docs).map(([path, doc]: [string, any]) => {
                        let slug = path.split("/")[2];
                        if (doc.version) slug = slug + "/" + doc.version;
                        if (doc.main) slug = slug + "/" + doc.main;
                        return (
                            <ListItem
                                key={slug}
                                onClick={() => navigate(`/${slug}`)}
                            >
                                <ListIcon>
                                    <Logo />
                                </ListIcon>
                                <ListInfo>
                                    <ListTitle>{doc.title}</ListTitle>
                                    <ListDescription>
                                        {doc.description}
                                    </ListDescription>
                                </ListInfo>
                            </ListItem>
                        );
                    })
                }
            </List>
        </div>
    );
}
