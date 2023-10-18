import "@testing-library/jest-dom";
import { render } from "@testing-library/react";
import Button from "./index.tsx";

test("it renders", () => {
    render(<Button type="colored" />);
});
