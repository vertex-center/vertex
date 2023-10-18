import "@testing-library/jest-dom";
import { render } from "@testing-library/react";
import MaterialIcon from "./index.tsx";

test("it renders", () => {
    render(<MaterialIcon name="deployed_code_update" />);
});
