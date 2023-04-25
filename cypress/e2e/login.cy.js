import { uid } from 'uid';

describe('Login page responds', () => {
    it('Checks if /login responds', () => {
        cy.visit('localhost:8080/login')
    })
})

describe('Login error messages', () => {    
    it('Checks for "invalid password" error message', () => {
        // ---------------- SETUP ---------------- //

        const username = uid()

        // register
        cy.register(username, "user@example.com", "password", "password")

        // login with wrong password
        cy.login(username, "wrong")

        // ---------------- TEST ---------------- //

        // check that correct error message shows up
        cy.contains("Error: Invalid password")
    })

    it('Checks for "invalid username" error message', () => {
        // ---------------- SETUP ---------------- //

        // login without registering
        cy.login(uid(), "password")

        // ---------------- TEST ---------------- //

        // check that correct error message shows up
        cy.contains("Error: Invalid Username")
    })
})

describe('Login user logs in', () => {    
    it('Checks if login logs the user in', () => {
        // ---------------- SETUP ---------------- //

        const username = uid()

        // register user
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        // login
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()

        // ---------------- TEST ---------------- //
        
        // check that the page contains text which appears for logged in users
        cy.url().should('contain', '/'+username)

        // ---------------- CLEANUP ---------------- //

        cy.get('a').contains('sign out').click()
    })
})