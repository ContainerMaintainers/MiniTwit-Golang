import { uid } from 'uid';

describe('Login page responds', () => {
    it('Checks if /login responds', () => {
        cy.visit('localhost:8080/login')
    })
})

describe('Login error messages', () => {    
    it('Checks for "invalid password" error message', () => {
        const username = uid()
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="password"]').type('wrong')
        cy.get('input').contains('Sign In').click()
        cy.contains("Error: Invalid password")
    })

    it('Checks for "invalid username" error message', () => {
        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(uid())
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        cy.contains("Error: Invalid Username")
    })
})

describe('Login user logs in', () => {    
    it('Checks if login logs the user in', () => {
        const username = uid()
        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.contains('You were logged in')

        cy.get('a').contains('sign out').click()
    })
})