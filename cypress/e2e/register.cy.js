import { uid } from 'uid';

describe('Register page responds', () => {
    it('Checks if /register responds', () => {
        cy.visit('localhost:8080/register')
    })
})

describe('Register error messages', () => {
    it('Checks for "all fields are required" error message', () => {
        cy.visit('localhost:8080/register')
        cy.get('input').contains('Sign Up').click()
        cy.contains("Error: All fields are required")
    })

    it('Checks for "you have to enter a valid email address" error message', () => {
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type('user')
        cy.get('input[name="email"]').type('example')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()
        cy.contains("Error: You have to enter a valid email address")
    })

    it('Checks for "the two passwords do not matchs" error message', () => {
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type('user')
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('wrong')
        cy.get('input').contains('Sign Up').click()
        cy.contains("Error: The two passwords do not match")
    })

    it('Checks for "the username is already taken" error message', () => {
        const username = uid()
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.contains("Error: The username is already taken")
    })
})

describe('Redirection', () => {
    it('Checks that register redirects to login page', () => {
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(uid())
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()
        cy.url().should('contain', '/login')
    })
})