describe('Check app is running', () => {
  it('Checks if /ping responds', () => {
    cy.request('localhost:8080/ping')
  })
})