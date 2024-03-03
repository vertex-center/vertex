import {
    ListDescription,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    Logo,
    MaterialIcon,
    Title,
} from "@vertex-center/components";
import styles from "./AllDocs.module.sass";
import { useNavigate } from "react-router-dom";

const docs = import.meta.glob("/{docs,api}/*/{doc,api}.json", {
    eager: true,
});

type DocsCategoryProps = {
    category: string;
    title: string;
};

function DocsCategory({ category, title }: DocsCategoryProps) {
    const navigate = useNavigate();

    return (
        <div className={styles.category}>
            <Title variant="h2">{title}</Title>
            <div className={styles.grid}>
                {Object.entries(docs).map(([path, doc]: [string, any]) => {
                    if (!path.startsWith("/" + category + "/")) return null;
                    let slug = path.split("/")[2];
                    if (category === "api") slug = "api-" + slug;
                    if (doc.version) slug = slug + "/" + doc.version;
                    if (doc.main) slug = slug + "/" + doc.main;
                    return (
                        <ListItem
                            key={slug}
                            onClick={() => navigate(`/${slug}`)}
                        >
                            <ListIcon>
                                {category === "docs" ? (
                                    <Logo />
                                ) : (
                                    <MaterialIcon icon="api" />
                                )}
                            </ListIcon>
                            <ListInfo>
                                <ListTitle>{doc.title}</ListTitle>
                                <ListDescription>
                                    {doc.description}
                                </ListDescription>
                            </ListInfo>
                        </ListItem>
                    );
                })}
            </div>
        </div>
    );
}

export default function AllDocs() {
    return (
        <div className={styles.all}>
            <DocsCategory category="docs" title="Documentations" />
            <DocsCategory category="api" title="Rest APIs" />
        </div>
    );
}
