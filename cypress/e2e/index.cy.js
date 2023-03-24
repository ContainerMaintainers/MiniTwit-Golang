describe('Main page responds', () => {
    it('Checks if index returns index page', () => {
        cy.visit('localhost:8080/')
    })
})

describe('Navigation buttons', () => {
    it('Checks for the existance of navigation buttons', () => {
        cy.visit('localhost:8080/')
        cy.get('a').contains('public timeline')
        cy.get('a').contains('sign up')
        cy.get('a').contains('sign in')
    })

    it('Checks that public timeline link leads to correct place', () => {
        cy.visit('localhost:8080/')
        cy.get('a').contains('public timeline').click()
        cy.url().should('contain', '/public')
    })

    it('Checks that sign up link leads to correct place', () => {
        cy.visit('localhost:8080/')
        cy.get('a').contains('sign up').click()
        cy.url().should('contain', '/register')
    })

    it('Checks that sign in link leads to correct place', () => {
        cy.visit('localhost:8080/')
        cy.get('a').contains('sign in').click()
        cy.url().should('contain', '/login')
    })

    it('Checks that messages exist', () => {
        cy.visit('localhost:8080/')
        cy.get('.messages > li').its('length').should('be.greaterThan', 2)
    })
})