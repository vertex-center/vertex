import {
    Code,
    Paragraph,
    Table,
    Title,
    useTitle,
} from "@vertex-center/components";
import { useLocation } from "react-router-dom";
import { Fragment, useMemo } from "react";

const docs = import.meta.glob("/api/*/api.json", {
    eager: true,
});

type ApiMethodsProps = {
    responses: any;
};

function ApiMethods(props: Readonly<ApiMethodsProps>) {
    const { responses } = props;
    return (
        <Table>
            <thead>
                <tr>
                    <th>Code</th>
                    <th>Description</th>
                </tr>
            </thead>
            <tbody>
                {Object.entries(responses ?? {}).map(([code, resp]) => {
                    return (
                        <tr key={code}>
                            <td>{code}</td>
                            <td>{resp.description}</td>
                        </tr>
                    );
                })}
            </tbody>
        </Table>
    );
}

type ApiProps = {
    api: any;
    tag: string;
};

export function Api(props: Readonly<ApiProps>) {
    const { api, tag } = props;

    const location = useLocation();

    const app = location.pathname.split("/")?.[1];
    const doc: any = docs?.[`/api/${app}/api.json`];

    useTitle(doc?.title ?? "-");

    const routes = useMemo(() => {
        if (api === undefined) return [];

        let routes: any = [];
        Object.entries(api.paths ?? {}).forEach(
            ([path, methods]: [string, any]) => {
                Object.entries(methods).forEach(
                    ([method, operation]: [string, any]) => {
                        if (operation.tags === undefined) return;
                        if (!operation.tags.includes(tag)) return;
                        routes.push({ path, method, operation });
                    }
                );
            }
        );
        return routes;
    }, [api, tag]);

    if (api === undefined) return null;

    return (
        <div className="documentation">
            {routes.map((route: any) => {
                console.log(route);
                const { path, method, operation } = route;
                const { operationId, summary, description, responses } =
                    operation;

                return (
                    <Fragment key={operationId}>
                        <Title variant="h2">{summary}</Title>
                        <Code language="bash">
                            {`${method.toUpperCase()} ${path}`}
                        </Code>
                        <Paragraph>{description}</Paragraph>
                        <div>
                            <ApiMethods responses={responses} />
                        </div>
                    </Fragment>
                );
            })}
        </div>
    );
}
