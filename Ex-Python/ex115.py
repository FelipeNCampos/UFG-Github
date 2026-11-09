# menu
import ex115m.m115
import ex115m.dds
arq = 'Registro.txt'
if ex115m.dds.arqexiste(arq):
    print('verdade')
else :
    ex115m.dds.criararq('Registro.txt')

print('---'*15)
print(f'{'MENU PRINCIPAL':^45}')
print('---'*15)
print(f'1 - Ver dados cadastrados.')
print(f'2 - Cadastrar uma nova pessoa.')
print(f'3 - Sair do programa.')
print('---'*15)
rsp = ex115m.m115.choisse('\033[33mSua opção : \033[m')