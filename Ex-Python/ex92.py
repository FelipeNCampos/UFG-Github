import datetime
pp={}
pp['Nome']=input('Digite seu nome :').strip().title()
nc = input('Digite sua data de nascimento :')
cdk = datetime.datetime.now()
date = cdk.date()
date
yr = date.strftime('%Y')
pp['Idade']=int(yr)-int(nc)
pp['CTPS']=input('Carteira de Trabalho (0 se não tem) :')
if int(pp['CTPS'])!=0:
    pp['Ano de contratação']=input(f'Digite o ano de contratação {pp['Nome']} :')
    pp['Salário']=input('Digite o salário :')
    pp['Idade de aposentadoria']=int(pp['Ano de contratação'])+30
print('-=-'*10)
for c in pp:
    print(f'{c} tem o valor de {pp[c]}')