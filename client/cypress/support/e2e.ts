import "./commands";

beforeEach(() => {
    cy.request("POST", "http://localhost:6130/api/hard-reset");
});
