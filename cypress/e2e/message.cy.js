import { uid } from 'uid';

describe('Create message', () => {
    
    it('Check if create message box exists', () => {
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
        
        cy.visit('localhost:8080/'+username)

        cy.contains('What\'s on your mind ?')
        cy.get('input').contains("Share")
        cy.get('input[name="text"]')

        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in private timeline', () => {
        const message = uid()
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
        
        cy.visit('localhost:8080/'+username)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.contains('li', username).contains(message)

        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in public timeline', () => {
        const message = uid()
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
        
        cy.visit('localhost:8080/'+username)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.visit('localhost:8080/public')

        cy.contains('li', username).contains(message)

        cy.get('a').contains('sign out').click()
    })

    it('Check if message is in users timeline', () => {
        const message = uid()
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
        
        cy.visit('localhost:8080/'+username)

        cy.get('input[name="text"]').type(message)
        cy.get('input').contains("Share").click()

        cy.get('a').contains('sign out').click()

        cy.visit('localhost:8080/'+username)

        cy.contains('li', username).contains(message)

    })

    it('Check if messages are in correct order', () => {
        const message1 = uid()
        const message2 = uid()
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
        
        cy.visit('localhost:8080/'+username)

        cy.get('input[name="text"]').type(message1)
        cy.get('input').contains("Share").click()

        cy.visit('localhost:8080/'+username)

        cy.get('input[name="text"]').type(message2)
        cy.get('input').contains("Share").click()

        cy.get('.messages > li').eq(0).should('contain', message2)
        cy.get('.messages > li').eq(1).should('contain', message1)

        cy.get('a').contains('sign out').click()

    })
})
