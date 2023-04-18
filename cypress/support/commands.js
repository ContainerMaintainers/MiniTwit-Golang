// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

Cypress.Commands.add('register', (username, email, password1, password2) => {
    cy.visit('localhost:8080/register')
    cy.get('input[name="username"]').type(username)
    cy.get('input[name="email"]').type(email)
    cy.get('input[name="password"]').type(password1)
    cy.get('input[name="password2"]').type(password2)
    cy.get('input').contains('Sign Up').click()
})

Cypress.Commands.add('login', (username, password) => {
    cy.visit('localhost:8080/login')
    cy.get('input[name="username"]').type(username)
    cy.get('input[name="password"]').type(password)
    cy.get('input').contains('Sign In').click()
})

Cypress.Commands.add('createMessage', (username, message) => {
    cy.visit('localhost:8080/' + username)
    cy.get('input[name="text"]').type(message)
    cy.get('input').contains("Share").click()
})