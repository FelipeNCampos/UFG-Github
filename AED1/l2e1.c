/*Em uma máquina de autoatendimento bancário (ATM), um cliente pode realizar saques em dinheiro. Considere que na máquina há cédulas de 2, 5, 10, 20, 50 e 100 reais e que há uma disponibilidade ilimitada de notas de cada um destes valores. Faça um programa em C para receber como entradas o saldo da conta do cliente e o valor desejado do saque e, caso o cliente tenha saldo suficiente e o valor do saque seja possível considerando as cédulas disponíveis, determinar quantas notas de cada valor o cliente deve receber e informar o novo saldo. Para determinar os valores e as quantidades de cédulas que deverão ser entregues ao cliente devem ser consideradas duas possibilidades: 1) com número mínimo de notas (dinheiro menos trocado); 2) com número máximo de notas (dinheiro mais trocado);

 

Considerações sobre as entradas e saídas e respectivas formatações:

·         O saldo da conta deve ser valor real com duas casas decimais

·         O valor do saque deve ser valor inteiro maior do que zero
*/
#include <stdio.h>
#include <stdlib.h>

float saldo;

int confere_valor(int saque){
    if(saque>saldo){
        return 1;
    }else if (saque<=0){
        return 2;
    }else if (saque == 1 || saque == 3 ){
        return 3;
    }
    return 0;
}
void realizar_saque(int saque){
    printf("RESULTADOS\n\n");
    printf("Valor atualizado do saldo (em reais): %.2f\n\n",(float)saldo-saque);

    printf("Opcao de saque 1 (menor quantidade possivel de cedulas):\n\n");
    int cont;
    cont = saque;
    int temp;
    temp = 0;
    while (cont>=100 && !(cont-100==3 || cont-100==1)){
        temp++;
        cont-=100;
    }
    
    if(temp!=0){
        printf("Cedulas de 100 reais: %i\n\n",temp);
        temp = 0;
    }

    while (cont>=50 && !(cont-50==3 || cont-50==1)){
        temp++;
        cont-=50;
    }
    if(temp!=0){
        printf("Cedulas de 50 reais: %i\n\n",temp);
        temp = 0;
    }

    while (cont>=20 && !(cont-20==3 || cont-20==1)){
        temp++;
        cont-=20;
    }
    if(temp!=0){
        printf("Cedulas de 20 reais: %i\n\n",temp);
        temp = 0;
    }

    while (cont>=10 && !(cont-10==3 || cont-10==1)){
        temp++;
        cont-=10;
    }
    if(temp!=0){
        printf("Cedulas de 10 reais: %i\n\n",temp);
        temp = 0;
    }

    while (cont>=5 && !(cont-5==3 || cont-5==1)){
        temp++;
        cont-=5;
    }
    if(temp!=0){
        printf("Cedulas de 5 reais: %i\n\n",temp);
        temp = 0;
    }

    while (cont>=2){
        temp++;
        cont-=2;
    }
    if(temp!=0){
        printf("Cedulas de 2 reais: %i\n\n",temp);
        temp = 0;
    }


    printf("Opcao de saque 2 (maior quantidade possivel de cedulas):\n\n");
    if(saque%2==0){
        printf("Cedulas de 2 reais: %d\n",saque/2);
    }else {
        printf("Cedulas de 2 reais: %d\n",(saque-5)/2);
        printf("Cedulas de 5 reais: %d\n",1);
    }
}



void main(){
    printf("Entre com o saldo da conta a ser sacada (em reais): ");
    scanf("%f",&saldo);printf("\n\n");

    int saque;
    int flag; 
    do{
        printf("Entre com o valor do saque (em reais): ");
        scanf("%i",&saque);printf("\n\n");
        flag = confere_valor(saque);
        switch (flag)
        {
            case 1:
                printf("Saldo insuficiente para valor desejado\n\n");
                break;
            case 2:
                printf("Valor do saque invalido\n\n");
                break;
            case 3:
                printf("Valor do saque incompativel com notas disponiveis\n\n");
                break;
        default:
            realizar_saque(saque);
            break;
        }
    }while(flag!=0);
}
