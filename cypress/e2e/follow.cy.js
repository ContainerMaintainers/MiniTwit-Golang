import { uid } from 'uid';

describe('Follow', () => {
    
    it('User follows and message shows up', () => {
        const message = uid()
        const username1 = uid()
        const username2 = uid()

        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.visit('localhost:8080/'+username1)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.get('a').contains('sign out').click()

        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.visit('localhost:8080/'+username1)

        cy.get('.follow').click()

        cy.visit('localhost:8080/')

        cy.contains('li', username1).contains(message)

        cy.get('a').contains('sign out').click()
    })

    it('User unfollows and message disappears', () => {
        const message = uid()
        const username1 = uid()
        const username2 = uid()

        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username1)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.visit('localhost:8080/'+username1)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.get('a').contains('sign out').click()

        cy.visit('localhost:8080/register')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="email"]').type('example@mail.com')
        cy.get('input[name="password"]').type('password')
        cy.get('input[name="password2"]').type('password')
        cy.get('input').contains('Sign Up').click()

        cy.visit('localhost:8080/login')
        cy.get('input[name="username"]').type(username2)
        cy.get('input[name="password"]').type('password')
        cy.get('input').contains('Sign In').click()
        
        cy.visit('localhost:8080/'+username1)

        cy.get('.follow').click()

        cy.visit('localhost:8080/')

        cy.contains('li', username1).contains(message)

        cy.visit('localhost:8080/'+username1)

        cy.get('.unfollow').click()

        cy.visit('localhost:8080/')

        cy.contains('li', username1).should('not.exist')

        cy.get('a').contains('sign out').click()
    })

})
