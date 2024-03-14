import {
    ListDescription,
    ListInfo,
    ListTitle,
    Logo,
    Title,
} from "@vertex-center/components";
import styles from "./AllDocs.module.sass";
import { useNavigate } from "react-router-dom";
import { Grid } from "@vertex-center/components/lib/components/Grid/Grid.tsx";
import { Card } from "@vertex-center/components/lib/components/Card/Card.tsx";
import { Vertical } from "@vertex-center/components/lib/components/Layout/Layout.tsx";
import { Graph } from "@phosphor-icons/react";

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
            <Grid className={styles.grid} columnSize={300}>
                {Object.entries(docs).map(([path, doc]: [string, any]) => {
                    if (!path.startsWith("/" + category + "/")) return null;
                    let slug = path.split("/")[2];
                    if (category === "api") slug = "api-" + slug;
                    if (doc.version) slug = slug + "/" + doc.version;
                    if (doc.main) slug = slug + "/" + doc.main;
                    return (
                        <Card key={slug} onClick={() => navigate(`/${slug}`)}>
                            <Vertical gap={20}>
                                <div>
                                    {category === "docs" ? (
                                        <Logo />
                                    ) : (
                                        <Graph size={32} />
                                    )}
                                </div>
                                <ListInfo>
                                    <ListTitle>{doc.title}</ListTitle>
                                    <ListDescription>
                                        {doc.description}
                                    </ListDescription>
                                </ListInfo>
                            </Vertical>
                        </Card>
                    );
                })}
            </Grid>
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
