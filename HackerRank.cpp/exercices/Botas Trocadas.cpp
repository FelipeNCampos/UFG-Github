#include<bits/stdc++.h>
using namespace std;



int nn,cont,i,tamanhotemp;
string templado;

int main(){
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    cin >> nn;
    vector<int>dir;
    vector<int>esq;
    for( i=0;i<nn;i++){
        cin>>tamanhotemp;
        cin>>templado;
        if (templado=="D"){
            dir.push_back(tamanhotemp);
        }else{
            esq.push_back(tamanhotemp);
        }

    }
    if (dir.size() > esq.size()){
        for (int c=0;c<dir.size();c++){
            for (int j=0;j<esq.size();j++){
                if (dir[c]==esq[j]){
                    cont++;
                    esq[j]=0;
                    break;
                }
            }
        }
    }else{
        for (int c=0;c<esq.size();c++){
            for (int j=0;j<dir.size();j++){
                if (esq[c]==dir[j]){
                    cont++;
                    dir[j]=0;
                    break;
                }
            }
        }
    }
    cout<<cont<<"\n";
    return 0;
}
